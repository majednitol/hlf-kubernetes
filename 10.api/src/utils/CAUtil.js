/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const adminUserId = 'admin';
const adminUserPasswd = 'adminpw';

/**
 *
 * @param {*} FabricCAServices
 * @param {*} ccp
 */
import { performance } from 'perf_hooks';

export async function registerAndEnrollUser(caClient, wallet, orgMspId, userId, affiliation, secret, encryptionKey) {
  

    try {
        
        const userIdentity = await wallet.get(userId);

        if (userIdentity) {
            console.log(`User ${userId} exists.`);
            return;
        }

        // Get admin identity from the wallet
    
        const adminIdentity = await wallet.get(adminUserId);
        if (!adminIdentity) {
            throw new Error('Admin identity not found.');
        }

        // Prepare admin context for registration
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, adminUserId);

        // Register the user
        const registerStart = performance.now();
        const secretData = await caClient.register({
            affiliation: affiliation,
            enrollmentID: userId,
            role: 'client'
        }, adminUser);

        // Enroll the user
        const enrollment = await caClient.enroll({
            enrollmentID: userId,
            enrollmentSecret: secretData
        });

        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: orgMspId,
            type: 'X.509'
        };

        await wallet.put(userId, x509Identity);


        // Log timing results
        console.log(`User ${userId} enrolled successfully.`);

    } catch (error) {
        console.error(`Error: ${error}`);
    }
}

export function buildCAClient(FabricCAServices, ccp, caHostName) {

	const caInfo = ccp.certificateAuthorities[caHostName];
	const caTLSCACerts = caInfo.tlsCACerts.pem;
	const caClient = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: true }, caInfo.caName);
	console.log(`Built a CA Client named ${caInfo.caName}`);
	return caClient;
}

export async function enrollAdmin(caClient, wallet, orgMspId) {
	try {
		// Check to see if we've already enrolled the admin user.
		const identity = await wallet.get(adminUserId);
		if (identity) {
			console.log('An identity for the admin user already exists in the wallet');
			return;
		}

		console.log("Admin Identity not found... Enroll admin")
		// Enroll the admin user, and import the new identity into the wallet.
		const enrollment = await caClient.enroll({ enrollmentID: adminUserId, enrollmentSecret: adminUserPasswd });
		const x509Identity = {
			credentials: {
				certificate: enrollment.certificate,
				privateKey: enrollment.key.toBytes(),
			},
			mspId: orgMspId,
			type: 'X.509',
		};

		await wallet.put(adminUserId, x509Identity);
		console.log('Successfully enrolled admin user and imported it into the wallet');
	} catch (error) {
		console.error(`Failed to enroll admin user : ${error}`);
	}
}
