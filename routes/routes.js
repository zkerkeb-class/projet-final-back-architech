import { Router } from "express";
import {
  getAllUsers,
  getUserById,
  createUser,
  updateUser,
  getUserByEmail,
  deleteUser,
  updateAccountMoney,
} from "../controllers/controller.js";

const router = Router();

router.get("/", getAllUsers);
router.get("/:id", getUserById);
router.get("/mail/:mail", getUserByEmail);
router.post("/", createUser);
router.put("/:id", updateUser);
router.delete("/:id", deleteUser);
router.patch("/:id/account-money", updateAccountMoney);

export default router;