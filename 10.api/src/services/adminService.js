import { smartContract } from "./smartContract.js";
export async function getAdmin(request) {
  try {
    const userId = request.userId
    console.log("userId", userId)
    const contract = await smartContract(request, userId)
    let result = await contract.evaluateTransaction("GetAdmin", userId);
    console.log("result", result)
    return JSON.parse(result);
  } catch (error) {
    console.log(error)
  }
}

export async function setAdmin(request) {
  try {
    let data = request.data;
    const userId = data.userId
    const contract = await smartContract(request, userId)
    let result = await contract.submitTransaction(
      "SetAdmin",
      data.userId,
      data.name,
      data.gender,
      data.birthday,
      data.emailAddress,
      data.location
    );
    console.log("Transaction Result:", result);

    return result;
  } catch (error) {
    console.error("Error in createAsset:", error);
    throw error;
  }
}
export async function allAdminData(request) {
  try {
    const userId = request.userId
    console.log("userId", userId)
    const contract = await smartContract(request, userId)
    let result = await contract.evaluateTransaction("GetAllAdmindata");
    console.log("result", result)
    return JSON.parse(result);
  } catch (error) {
    console.log(error)
  }
}

export async function giveConfirmation(request) {
  try {
    
    const adminId = request.adminId
    const contract = await smartContract(request, adminId)
    let result = await contract.submitTransaction(
      "GiveConfirmation",
      request.userId,
      adminId
    );
    console.log("Transaction Result:", result);

    return result;
  } catch (error) {
    console.error("Error in createAsset:", error);
    throw error;
  }
}