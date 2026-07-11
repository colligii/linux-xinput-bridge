const fs = require('fs');

class Controller {
    constructor(DEVICE, config) {
        this.controllerButtons = {};
        this.joystickDecodeAlgo = {};
        this.btnMapping = config?.btnMapping;
        this.groups = config?.joystickGroup ?? [];
        this.composedKeys = config?.composedKeys;

        this.codeToGroup = {};
        this.group = {};

        for(let item of Object.values(this.btnMapping)) {
            this.controllerButtons[item] = 0;
        }

        for(let [key, value, paramsMapping] of config.joystickDecodeAlgo) {
            this.joystickDecodeAlgo[key] = [
                require(value),
                paramsMapping
            ];
        }

        for(let [key, value] of this.groups) {
            const json = {};

            for(let btnCode of value) {
                delete this.controllerButtons[this.btnMapping[btnCode]]
                json[this.btnMapping[btnCode]] = 0;
                this.codeToGroup[this.btnMapping[btnCode]] = key;
            }

            this.controllerButtons[key] = JSON.parse(JSON.stringify(json));

            this.group[key] = JSON.parse(JSON.stringify(json))
        }


        // let biggestValue = 0;

        this.tryToReadFileSync(DEVICE)
            .then(() => {
                this.startControllerStream(DEVICE)
            })
            .catch((e) => {
                if(e?.message?.includes("EACCES: permission denied, access ")) {
                    throw new Error("You may should change the defaultConfig.json and then run (npm run givePermission) or (sudo ./givePermission.sh) before run npm start");
                }
                throw new Error(e)
            })

    }

    startControllerStream(DEVICE) {

        const stream = fs.createReadStream(DEVICE);



        stream.on('data', (buffer) => {
            for (let offset = 0; offset < buffer.length; offset += 24) {

                const type = buffer.readUInt16LE(offset + 16);
                let _code = buffer.readUInt16LE(offset + 18);
                const _value = buffer.readInt32LE(offset + 20);

                if(this.composedKeys[`${_code};${type}`]) {
                    _code = `${_code};${type}`;
                }

                const { code, value } = this.getFieldValue(_code, _value)


                        // console.log({
                        //     type,
                        //     code: _code,
                        //     value: _value
                        // })


                if(this.cb && code) {
                    
                    
                    const groupCode = this.codeToGroup[code];
                    const isCodeToGroup = !!groupCode;
                    let _controllersButtons = JSON.parse(JSON.stringify(this.controllerButtons));
                    if(isCodeToGroup) {
                        this.group[groupCode][code] = value;

                        if(this.joystickDecodeAlgo?.[groupCode]?.[0]) {
                            const x_key = this.joystickDecodeAlgo[groupCode][1].x;
                            const y_key = this.joystickDecodeAlgo[groupCode][1].y;
                            const {x, y} = this.joystickDecodeAlgo[groupCode][0](
                                this.group[groupCode][x_key],
                                this.group[groupCode][y_key],
                            );

                            _controllersButtons[groupCode] = {
                                [x_key]: x,
                                [y_key]: y
                            }
                            
                        } else {
                            _controllersButtons[groupCode] = this.group[groupCode]
                        }

                        // biggestValue = biggestValue < value ? value : biggestValue;
                    } else {
                        _controllersButtons[code] = _value;
                    }

                    if(JSON.stringify(_controllersButtons) !== JSON.stringify(this.controllerButtons)) {
                        this.controllerButtons = JSON.parse(JSON.stringify(_controllersButtons))
                        this.cb(this.controllerButtons)
                    }

                }
                // else 
                    // console.log({
                    //     type,
                    //     code: _code,
                    //     value: _value
                    // })
            }
        });
    }

    tryToReadFileSync(path) {
        return new Promise((resolve, reject) => {

            fs.access(path, fs.constants.R_OK, (err) => {
                if (err) {
                    reject(err)
                } else {
                    resolve(true);
                }
            });
        })
    }

    getFieldValue(_code, value) {
        let code = this.btnMapping?.[String(_code)];
        
        return {
            code,
            value
        }; 
    }

    setCb(cb) {
        this.cb = cb;
        this.cb(this.controllerButtons)
    }
}

module.exports = Controller