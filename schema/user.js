import mongoose from "mongoose";

const userSchema = new mongoose.Schema({
    id: {
        type: Number,
        required: true,
        unique: true,
    },
    username: {
        type: String,
        required: true,
    },
    account_money: {
       type: Number,
    },
});

export default mongoose.model("user", userSchema);
