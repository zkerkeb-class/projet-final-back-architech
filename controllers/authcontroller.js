import crypto from "crypto";
import jwt from "jsonwebtoken";
import User from "../schemas/user.js";
import MagicToken from "../schemas/magicLink.js";
import { sendMagicLink } from "../utils/sendEmail.js";

// POST /api/auth/login
// User requests a magic link
export const requestMagicLink = async (req, res) => {
  try {
    const { email } = req.body;
    if (!email) return res.status(400).json({ message: "Email is required" });

    // Create user if they don't exist yet
    let user = await User.findOne({ email });
    if (!user) {
      user = new User({ email, username: email.split("@")[0] });
      await user.save();
    }

    // Delete any existing tokens for this email
    await MagicToken.deleteMany({ email });

    // Generate a secure random token
    const token = crypto.randomBytes(32).toString("hex");

    // Save token to database, expires in 15 minutes
    await MagicToken.create({
      email,
      token,
      expiresAt: new Date(Date.now() + 15 * 60 * 1000),
    });

    // Send the magic link email
    await sendMagicLink(email, token);

    res.status(200).json({ message: "Magic link sent to your email!" });
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};

// GET /api/auth/verify?token=xxxxx
// User clicks the magic link
export const verifyMagicLink = async (req, res) => {
  try {
    const { token } = req.query;
    if (!token) return res.status(400).json({ message: "Token is required" });

    // Find the token in the database
    const magicToken = await MagicToken.findOne({ token });
    if (!magicToken) return res.status(400).json({ message: "Invalid or expired token" });

    // Check if token has expired
    if (magicToken.expiresAt < new Date()) {
      await MagicToken.deleteOne({ token });
      return res.status(400).json({ message: "Token has expired" });
    }

    // Find the user
    const user = await User.findOne({ email: magicToken.email });
    if (!user) return res.status(404).json({ message: "User not found" });

    // Delete the used token (one time use only)
    await MagicToken.deleteOne({ token });

    // Generate a JWT session token
    const jwtToken = jwt.sign(
      { id: user._id, email: user.email },
      process.env.JWT_SECRET,
      { expiresIn: "7d" }
    );

    res.status(200).json({
      message: "Login successful!",
      token: jwtToken,
      user: {
        id: user._id,
        email: user.email,
        username: user.username,
      },
    });
  } catch (error) {
    res.status(500).json({ message: error.message });
  }
};