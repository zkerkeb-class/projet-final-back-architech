import nodemailer from "nodemailer";

const transporter = nodemailer.createTransport({
  host: "smtp.gmail.com",
  port: 465,
  secure: true,
  auth: {
    user: process.env.EMAIL_USER,
    pass: process.env.EMAIL_PASS,
  },
});

export const sendMagicLink = async (email, token) => {
  // Use the server URL for the verification link
  //J'ai enleve /api/ dans /auth/verify?token=${token}
  const magicLink = `${process.env.APP_URL}/auth/verify?token=${token}`;

  await transporter.sendMail({
    from: process.env.EMAIL_USER,
    to: email,
    subject: "Lien de connexion Bayment",
    html: `
      <div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; color: #1e293b;">
        <h2 style="color: #3b82f6;">Connectez-vous à Bayment</h2>
        <p>Cliquez sur le bouton ci-dessous pour vous connecter. Expiration dans 15 minutes.</p>
        <div style="text-align: center; margin: 30px 0;">
          <a href="${magicLink}" style="background-color: #3b82f6; color: white; padding: 14px 28px; border-radius: 8px; text-decoration: none; font-weight: bold; display: inline-block;">Me connecter</a>
        </div>
        <p style="color: #64748b; font-size: 14px;">Si vous n'êtes pas à l'origine de cette demande, vous pouvez ignorer cet e-mail.</p>
      </div>
    `,
  });
};