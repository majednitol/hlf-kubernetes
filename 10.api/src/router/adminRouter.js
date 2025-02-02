
import express from 'express' 
import { adminData, confirmation, createAdminAccount, getAllAdmin, getAllAdmindata, shareDataByAdmin } from '../controllers/adminController.js';
import authenticate from '../middleware/authenticate.js';

const adminRouter = express.Router()
adminRouter.post("/create-admin-account", createAdminAccount)
adminRouter.get('/getAdminData',authenticate, adminData);
adminRouter.get('/getAllAdminData',authenticate, getAllAdmindata);
adminRouter.get('/getAllAdmin', getAllAdmin);
adminRouter.post("/give-confirmation", authenticate, confirmation)
adminRouter.post("/share-data-by-admin",authenticate,shareDataByAdmin)

export default adminRouter