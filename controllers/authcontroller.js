import crypto from "crypto";
import jwt from "jsonwebtoken";
import User from "../schemas/user.js";
import MagicToken from "../schemas/magicLink.js";
import { sendMagicLink } from "../utils/sendEmail.js";

// POST /auth/login
// User requests a magic link
export const requestMagicLink = async (req, res) => {
  try {
    const { mail } = req.body;
    if (!mail) return res.status(400).json({ message: "Email is required" });

    // Create user if they don't exist yet
    let user = await User.findOne({ mail });
    if (!user) {
      user = new User({ mail, username: mail.split("@")[0] });
      await user.save();
    }

    // Delete any existing tokens for this email
    await MagicToken.deleteMany({ mail });

    // Generate a secure random token
    const token = crypto.randomBytes(32).toString("hex");

    // Save token to database, expires in 15 minutes
    await MagicToken.create({
      mail,
      token,
      expiresAt: new Date(Date.now() + 15 * 60 * 1000),
    });

    // Send the magic link email
    await sendMagicLink(mail, token);

    res.status(200).json({ message: "Magic link sent to your email!" });
  } catch (error) {
    console.error("Error in requestMagicLink:", error);
    res.status(500).json({ message: error.message });
  }
};

// GET /auth/verify?token=xxxxx
// User clicks the magic
export const verifyMagicLink = async (req, res) => {
  try {
    const { token } = req.query;
    if (!token) return res.status(400).send("<h1>Token is required</h1>");

    // Find the token in the database
    const magicToken = await MagicToken.findOne({ token });
    if (!magicToken) return res.status(400).send("<h1>Invalid or expired token</h1>");

    // Check if token has expired
    if (magicToken.expiresAt < new Date()) {
      await MagicToken.deleteOne({ token });
      return res.status(400).send("<h1>Token has expired</h1>");
    }

    // Find the user
    const user = await User.findOne({ mail: magicToken.mail });
    if (!user) return res.status(404).send("<h1>User not found</h1>");

    // Delete the used token (one time use only)
    await MagicToken.deleteOne({ token });

    // Generate a JWT session token
    const jwtToken = jwt.sign(
      { id: user._id, mail: user.mail },
      process.env.JWT_SECRET,
      { expiresIn: "7d" }
    );

    /*
     res.status(200).json({
     message: "Login successful!",
      token: jwtToken,
     user: {
       id: user._id,
      mail: user.mail,
        username: user.username,
      },
     });
    */
    // Redirect to the mobile app using its custom scheme
    // The scheme "bayment://" is defined in app.json
    const redirectUrl = `bayment://auth?token=${jwtToken}`;

    // Return a simple HTML page that tries to redirect or provides a button
    res.send(`
      <html>
        <body style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100vh; font-family: sans-serif; background-color: #0f172a; color: white;">
          <h1>Connexion réussie !</h1>
          <p>Vous allez être redirigé vers l'application...</p>
          <a href="${redirectUrl}" style="background-color: #3b82f6; color: white; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: bold; margin-top: 20px;">
            Ouvrir lapplication Bayment
          </a>
          <script>
            window.location.href = "${redirectUrl}";
          </script>
        </body>
      </html>
    `);
  } catch (error) {
    res.status(500).send('<h1>Erreur: ${error.message}</h1>');
  }
};