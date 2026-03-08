import { Router } from "express";
import { requestMagicLink, verifyMagicLink } from "../controllers/authcontroller.js";

const router = Router();

router.post("/login", requestMagicLink);      // User submits email
router.get("/verify", verifyMagicLink);       // User clicks magic link

export default router;