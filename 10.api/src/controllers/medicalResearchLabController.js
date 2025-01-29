import { getMedicalResearchLab, setMedicalResearchLab } from "../services/medicalResearchLabService.js";
const chaincodeName = "basic";
const channelName = "mychannel"
export async function createMedicalResearchLabAccount(req, res) {
    try {
        let payload = {
            "org": req.body.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "data": req.body.data
        }
        console.log("payload", payload)
        let result = await setMedicalResearchLab(payload);
        console.log(result)
        res.send(result)
    } catch (error) {
        console.log(error)
        res.status(500).send(error)
    }
}

export async function medicalResearchLabData(req, res) {
    try {
        let payload = {
            "org": req.org,
            "channelName": channelName,
            "chaincodeName": chaincodeName,
            "labID": req.query.labID ? req.query.labID : req.userId
        }
        console.log("payload", payload)
        let result = await getMedicalResearchLab(payload);
        console.log("result app", result)
        res.json(result)
    } catch (error) {
        console.log(error)
        res.status(500).send(error)
    }
}