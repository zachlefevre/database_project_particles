console.log("vector.js loaded");

class Pvector{
    constructor(x_, y_){
        this.x = x_;
        this.y = y_;
    }
    add(pToAdd){
        this.x += pToAdd.x;
        this.y += pToAdd.y
    }
    sub(pToSub){
        this.x -= pToSub.x;
        this.y -= pToSub.y;
    }
    multiply(numToMult){
        this.x *= numToMult;
        this.y *= numToMult; 
    }
    divide(numToDiv){
        this.x /= numToDiv;
        this.y /= numToDiv;
    }
    magnitude(){
        var x = sqrt(abs((this.x * this.x)) + abs((this.y * this.y)));
        return x
    }
    normalize(){
        var mag = this.magnitude();
        if(mag != 0){
            this.divide(mag);
        }
    }
    dist(other){
        return Math.sqrt(Math.pow(this.x - other.x, 2) + Math.pow(this.y - other.y, 2))
    }
    get(){
        return(new Pvector(this.x, this.y));
    }


}
function subVectors(PvectorM, PvectorN){
    return new Pvector(PvectorM.x-PvectorN.x, PvectorM.y-PvectorN.y);
}