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
  const magicLink = `${process.env.APP_URL}/api/auth/verify?token=${token}`;

  await transporter.sendMail({
    from: process.env.EMAIL_USER,
    to: email,
    subject: "Your Magic Login Link",
    html: `
      <h2>Login to Bayment</h2>
      <p>Click the link below to log in. This link expires in 15 minutes.</p>
      <a href="${magicLink}">Click here to login</a>
      <p>If you didn't request this, ignore this email.</p>
    `,
  });
};