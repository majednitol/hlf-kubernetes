
import express from 'express' 
import { AddPatientPersonalData, createPatientAccount, patientData, patientDataFromDoctor, shareDataByPatient } from '../controllers/patientController.js';
import authenticate from '../middleware/authenticate.js';
const patientRouter = express.Router()
patientRouter.post("/create-patient-account", createPatientAccount)
patientRouter.get('/getPatientData',authenticate, patientData);
patientRouter.post('/add-patient-personal-data', authenticate,AddPatientPersonalData);
patientRouter.get("/data-from-doctor", authenticate, patientDataFromDoctor)
patientRouter.post("/share-data-by-patient",authenticate,shareDataByPatient)
export default patientRouter