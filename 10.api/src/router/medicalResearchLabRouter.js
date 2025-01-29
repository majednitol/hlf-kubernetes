
import express from 'express' 
import { createMedicalResearchLabAccount, medicalResearchLabData } from '../controllers/medicalResearchLabController.js';
import authenticate from '../middleware/authenticate.js';

const medicalResearchLabRouter = express.Router()
medicalResearchLabRouter.post("/create-medicalResearchLab-account", createMedicalResearchLabAccount)
medicalResearchLabRouter.get('/getMedicalResearchLabRouterData', authenticate,medicalResearchLabData);

export default medicalResearchLabRouter