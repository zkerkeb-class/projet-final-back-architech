import mongoose from "mongoose";

const userSchema = new mongoose.Schema({
    id: {
        type: Number,
        required: false,
        unique: true,
    },
    username: {
        type: String,
        required: false,
    },
    mail: {
        type: String,
        required: true,
        unique: true,
        lowercase: true,      
        trim: true,           
        match: [/^\S+@\S+\.\S+$/, "Please enter a valid email address"],

    },
    account_money: {
       type: Number,
    },
}, { versionKey: false });

export default mongoose.model("user", userSchema);
