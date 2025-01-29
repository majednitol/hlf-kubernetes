
import express from 'express'
import { addTopMedicine, createPharmacyCompanyAccount, pharmacyCompanyData } from '../controllers/pharmacyCompanyController.js';

import authenticate from '../middleware/authenticate.js';
const pharmacyCompanyRouter = express.Router()
pharmacyCompanyRouter.post("/create-pharmacyCompany-account", createPharmacyCompanyAccount)
pharmacyCompanyRouter.post("/add-top-medicine", authenticate, addTopMedicine)
pharmacyCompanyRouter.get('/getPharmacyCompanyData', authenticate, pharmacyCompanyData);

export default pharmacyCompanyRouter