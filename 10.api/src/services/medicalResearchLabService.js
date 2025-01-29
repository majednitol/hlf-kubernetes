
import { smartContract } from "./smartContract.js";
export async function getMedicalResearchLab(request) {
    try {
        const labID = request.labID
        console.log("labID", labID)
        const contract = await smartContract(request, labID)
        let result = await contract.evaluateTransaction("GetMedicalResearchLab", labID);
        console.log("result", result)
        return JSON.parse(result);
    } catch (error) {
        console.log(error)
    }
}

export async function setMedicalResearchLab(request) {
    try {
        let data = request.data;
        const labID = data.labID
        const contract = await smartContract(request, labID)
        let result = await contract.submitTransaction(
            "SetMedicalResearchLab",
            data.labID,
            data.name,
            data.licenseID,
            data.researchArea,
            data.labRating,
            data.emailAddress
        );
        console.log("Transaction Result:", result);

        return result;
    } catch (error) {
        console.error("Error in createAsset:", error);
        throw error;
    }
}