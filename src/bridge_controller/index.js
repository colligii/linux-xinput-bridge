const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs')

class BridgeController {
    initialize() {
        return new Promise((resolve, reject) => {
            this.goBinaryPath = path.join(process.cwd(), 'lib_go', 'xbox-driver');
            this.goProcess = spawn(this.goBinaryPath, [], {
                stdio: ['pipe', 'inherit', 'inherit']
            });

            fs.writeFileSync(path.join(process.cwd(), 'logs.txt'), ("Go driver PID:"+ this.goProcess.pid), 'utf-8')
    
            this.goProcess.on('error', (err) => {
                reject('Falha ao iniciar o driver em Go:', err);
            });
    
            this.goProcess.on('exit', (code) => {
                process.exit(code);
                reject(`Driver Go fechou com código: ${code}. Você rodou com sudo?`);
            });

            process.on('SIGINT', () => {
                this.goProcess.kill();
                process.exit();
            });
            
            resolve("Conectado ao driver nativo em Go!");
        })
    }

    sendEvent = (controllerState) => {
        if (!this.goProcess.stdin || !this.goProcess.stdin.writable) {
            throw new Error("O canal de comunicação com o Go não está pronto ou foi fechado.");
        }

        try {
            const jsonStr = JSON.stringify(controllerState) + '\n';
            
            this.goProcess.stdin.cork();
            
            this.goProcess.stdin.write(jsonStr, 'utf8');
            
            this.goProcess.stdin.uncork();
        } catch (err) {
            console.error(err)
            console.error("Erro ao enviar dados para o processo Go:", err.message);
        }
    }
}


module.exports = BridgeController