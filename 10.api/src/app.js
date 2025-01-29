import express from "express";
const app = express();
import { GetAssetHistory } from './query.js';
import { registerUser } from "./registerUser.js";
;
import cors from 'cors';
import doctorRouter from "./router/doctorRouter.js";
import patientRouter from "./router/patientRouter.js";
import adminRouter from "./router/adminRouter.js";
import pathologistRouter from "./router/pathologistRouter.js";
import medicalResearchLabRouter from "./router/medicalResearchLabRouter.js";
import pharmacyCompanyRouter from "./router/pharmacyCompanyRouter.js";
import prescriptionRouter from "./router/prescriptionRouter.js";
import profilePicRouter from "./router/profilePicRouter.js";
import userRouter from "./router/userRouter.js";
import dataRouter from "./router/dataRouter.js";
import gobalErrorHander from "./middleware/gobalErrorHander.js";
const chaincodeName = "basic";
const channelName = "mychannel"
app.use(cors())
app.use(express.json());

app.listen(4000, () => {
    console.log("server started");

})
app.use('/doctor', doctorRouter);
app.use("/patient", patientRouter)
app.use('/admin', adminRouter)
app.use("/pathologist", pathologistRouter)
app.use("/medicalResearchLab", medicalResearchLabRouter)
app.use("/pharmacyCompany", pharmacyCompanyRouter)
app.use("/prescription", prescriptionRouter)
app.use("/profilePic", profilePicRouter)
app.use("/user", userRouter)
app.use("/data", dataRouter)
app.use(gobalErrorHander)







app.post("/updateAsset", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await updateAsset(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})


app.post("/transferAsset", async (req, res) => {

    try {

        let payload = {
            "org": req.body.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await TransferAsset(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})


app.post("/deleteAsset", async (req, res) => {
    try {
        let payload = {
            "org": req.body.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await deleteAsset(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})




app.get('/getAssetHistory', async (req, res) => {
    try {
        let payload = {
            "org": req.query.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "userId": req.query.userId,
            "data": {
                id: req.query.id
            }
        }

        let result = await GetAssetHistory(payload);
        res.json(result)
    } catch (error) {
        res.status(500).send(error)
    }

});


