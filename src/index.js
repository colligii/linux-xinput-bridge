const fs = require('fs');
const path = require('path');
const Controller = require('./controller');
const BridgeController = require('./bridge_controller');

const controller = process.env.CONTROLLER;

async function readConfigControllerJson(controller) {
    try {
        return JSON.parse(
            fs.readFileSync(
                path.join(process.cwd(), 'config', `${controller}.json`)
            )
        )
    } catch(e) {
        return null
    }
}

async function readDefaultConfigJson() {
    try {
        return JSON.parse(
            fs.readFileSync(
                path.join(process.cwd(), `defaultConfig.json`)
            )
        )
    } catch(e) {
        return null
    }
}

async function main() {
    const defaultConfig = await readDefaultConfigJson();
    const config = await readConfigControllerJson(defaultConfig?.controller ?? controller)

    if(!config) {
        throw new Error("You should have an controller configuration file to continue");
    }

    const controller = new Controller(defaultConfig?.evtest, config)
    const bridgeController = new BridgeController();

    await bridgeController.initialize();

    controller.setCb((json) => {
        // console.log(json)
        bridgeController.sendEvent.bind(bridgeController)((json))
    })

}

main()