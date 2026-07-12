let biggestValue = 65535;
let middle = 32767;
let start = 0;
let deadzone = 5000;

function translateValues(x, y) {
    if((x > middle && x < biggestValue) && (y > start && y < middle)) {

        x = (x - biggestValue);
        // console.log(x, y)
    } 
    else if((x > start && x < middle) && (y > start && y < middle)) {
        // console.log('top-right', x, y)  
    } 
    else if((x > start && x < middle) && (y > middle && y < biggestValue)) {
        y = (y - biggestValue);
        
    } else if((x > middle && x < biggestValue) && (y > middle && y < biggestValue)) {
        x = (x - biggestValue);
        y = (y - biggestValue);
        
    }

    if(x < deadzone && x > deadzone * -1) {
        x = 0;
    }


    if(y < deadzone && y > deadzone * -1) {
        y = 0;
    }
    
    y = y * -1;

    return {x,y}
}  

module.exports = translateValues;