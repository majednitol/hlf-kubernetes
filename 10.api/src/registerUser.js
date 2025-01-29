import { Wallets } from "fabric-network";
import FabricCAServices from 'fabric-ca-client';

import { buildCAClient, registerAndEnrollUser, enrollAdmin, userExist } from "./utils/CAUtil.js";
import { buildWallet } from "./utils/AppUtils.js";

import { resolve } from 'path';
import { Utils as utils } from 'fabric-common';
import { getCCP } from "./common/buildCCP.js";
import createHttpError from "http-errors";
let config = utils.getConfig()
config.file(resolve('config.json'))
let walletPath;
export async function registerUser({ OrgMSP, userId, secret, encryptionKey }) {

    let org = Number(OrgMSP.match(/\d/g).join(""));
    walletPath = resolve("wallet")
    let ccp = getCCP(org)
    const caClient = buildCAClient(FabricCAServices, ccp, `ca-org${org}`);
    const wallet = await buildWallet(Wallets, walletPath);
    console.log("wallet ", wallet)
    await enrollAdmin(caClient, wallet, OrgMSP);
    await registerAndEnrollUser(caClient, wallet, OrgMSP, userId, `org${org}.department1`, secret, encryptionKey);
    return {
        wallet
    }
}



const _userExist = async ({ OrgMSP, userId }) => {
    let org = Number(OrgMSP.match(/\d/g).join(""));
    let ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const result = await userExist(wallet, userId);
    return result;
};
export { _userExist as userExist };
