import express from "express";
import "dotenv/config";
import connect from "./connect.js";
import userRouter from "./routes/routes.js";
import authRouter from "./routes/authroutes.js";

const app = express();

app.use(express.json());
app.use("/users", userRouter);
app.use("/auth", authRouter);

connect().then(() => {
  app.listen(process.env.PORT || 3000, () => {
    console.log(`Server running on port ${process.env.PORT || 3000}`);
  });
});