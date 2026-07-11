let biggestValue = 65535;
let middle = 32767;
let start = 0;
function translateValues(x, y) {
    if((x > middle && x < biggestValue) && (y > start && y < middle)) {
        y = y - middle;
        x = x - middle
    } 
    else if((x > start && x < middle) && (y > start && y < middle)) {
        x = (x - middle) * -1
        y = (y - middle) * -1
        
    } 
    else if((x > start && x < middle) && (y > middle && y < biggestValue)) {
        x = x - middle;
        y = y - middle;
    } else if((x > middle && x < biggestValue) && (y > middle && y < biggestValue)) {
        x = (x - middle) * -1;
        y = (y - middle) * -1;
    }

    x += middle;
    y += middle;

    return {x,y}
}  

module.exports = translateValues;