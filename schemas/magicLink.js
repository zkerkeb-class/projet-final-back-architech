import mongoose from "mongoose";

const magicLinkSchema = new mongoose.Schema({
  mail: {
    type: String,
    required: true,
  },
  token: {
    type: String,
    required: true,
    unique: true,
  },
  expiresAt: {
    type: Date,
    required: true,
  },
});

magicLinkSchema.index({ expiresAt: 1 }, { expireAfterSeconds: 0 });

export default mongoose.model("MagicLink", magicLinkSchema);