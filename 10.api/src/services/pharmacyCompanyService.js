
import { smartContract } from "./smartContract.js";
export async function getPharmacyCompany(request) {
    try {
        const companyID = request.companyID
        console.log("companyID", companyID)
        const contract = await smartContract(request, companyID)
        let result = await contract.evaluateTransaction("GetPharmacyCompany", companyID);
        console.log("result", result)
        return JSON.parse(result);
    } catch (error) {
        console.log(error)
    }
}

export async function setPharmacyCompany(request) {
    try {
        let data = request.data;
        const companyID = data.companyID
        const contract = await smartContract(request, companyID)
        let result = await contract.submitTransaction(
            "SetPharmacyCompany",
            data.companyID,
            data.name,
            data.licenseID,
            data.productInformation,
            data.pharmacyRating,
            data.emailAddress
        );
        console.log("Transaction Result:", result);

        return result;
    } catch (error) {
        console.error("Error in createAsset:", error);
        throw error;
    }
}

export async function addingTopMedicine(request) {
    try {
       
        const companyID = request.companyID
        console.log("companyID",companyID)
        const contract = await smartContract(request, companyID)
        let result = await contract.submitTransaction(
            "AddTopMedicine",
            request.companyID,
            request.medicine
        );
        console.log("Transaction Result:", result);

        return result;
    } catch (error) {
        console.error("Error in createAsset:", error);
        throw error;
    }
}