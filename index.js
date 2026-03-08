import express from "express";
import connect from "./connect.js";
import userRouter from "./routes/user.routes.js";
import authRouter from "./routes/auth.routes.js";

const app = express();

app.use(express.json());
app.use("/api/users", userRouter);
app.use("/api/auth", authRouter);

connect().then(() => {
  app.listen(process.env.PORT || 3000, () => {
    console.log(`Server running on port ${process.env.PORT || 3000}`);
  });
});
