import moongose from 'mongoose';

const connectDB = async () => {
    try {
        await moongose.connect("mongodb://localhost:27017/bluemint");
        console.log('Connected to Mongo DB succesfully');
    } catch (error) {
        console.error('Error connecting to MongoDB:', error);
    }
};

connectDB()