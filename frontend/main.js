console.log("main.js loaded");
var apiURL = "http://localhost:3080/api"
axios.defaults.baseURL = "http://localhost"
axios.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8';
axios.defaults.headers.post['crossDomain'] = true;

var maxBalls = 20
var ballToSizeScalingRatio = 3
function wallEvent(particle, wall, epoch, timestep) {
    console.log(particle.name + ' collided with wall ' + wall)
    data = {
        p: particle.name,
        wall: wall,
        epoch: epoch,
        timestep: timestep,
        mode: 'no-cors'
    }
    axios.post(apiURL + "/collision/wall", data).then(ret => console.log(ret)).catch(err => console.log(err, JSON.stringify(data)))
}
function ballEvent(particle1, particle2) {
    particle1.randomizeColor()
    particle2.randomizeColor()
    console.log(particle1, particle2)
}

class Liquid {
    constructor({ x_, y_, w_, h_, c_ }) {
        this.x = x_;
        this.y = y_;
        this.w = w_;
        this.h = h_;
        this.c = c_;
        this.color = 175;
    }
    display() {
        noStroke();
        fill(this.color);
        ellipseMode(CENTER)
        rect(this.x, this.y, this.w, this.h);
    }
}

class Mover {
    constructor({ x_, y_, vx, vy, ax, ay, mass_, topSpeed, name_ }) {
        this.name = name_
        this.location = new Pvector(x_, y_);
        this.velocity = new Pvector(vx, vy);
        this.acceleration = new Pvector(ax, ay);
        this.mass = mass_;
        this.topSpeed = topSpeed;
        this.w = this.mass * ballToSizeScalingRatio;
        this.h = this.mass * ballToSizeScalingRatio;
        this.randomizeColor()
    }
    display() {
        stroke(0);
        fill(this.r, this.g, this.b);
        ellipse(this.location.x, this.location.y, this.w, this.h);
    }
    randomizeColor() {
        this.r = random(0, 255)
        this.g = random(0, 255)
        this.b = random(0, 255)
    }
    update() {
        this.velocity.add(this.acceleration);
        this.location.add(this.velocity);
        this.acceleration.multiply(0);

    }
    checkForWalls() {
        if (this.location.x > width) {
            this.velocity.x *= -1;
            wallEvent(this, 'RIGHT', epoch, ts)
        } else if (this.location.x < 0) {
            this.velocity.x *= -1;
            wallEvent(this, 'LEFT', epoch, ts)
        }

        if (this.location.y > height) {
            this.velocity.y *= -1;
            wallEvent(this, 'BOTTOM', epoch, ts)
        } else if (this.location.y < 0) {
            this.velocity.y *= -1;
            wallEvent(this, 'TOP', epoch, ts)
        }
    }
    checkForCollision(ballArr) {
        for (var i = 0; i < ballArr.length; i++) {
            const ball = ballArr[i]
            if (this.name == ball.name) break;
            if (this.location.dist(ball.location) < this.w / 2 + ball.w / 2) {
                ballEvent(this, ball)
            }
        }
    }
    limit(max) {
        if (this.velocity.magnitude() > this.topSpeed) {
            this.velocity.normalize();
            this.velocity.multiply(this.topSpeed);
        }
    }
    applyForce(forceToAdd) {
        var force = forceToAdd.get();
        force.divide(this.mass);
        this.acceleration.add(force);

    }
    isInside(fluid) {
        if (this.location.x > fluid.x && this.location.x < (fluid.x + fluid.w) && this.location.y > fluid.y && this.location.y < (fluid.y + fluid.h)) {
            return true;
        } else {
            return false;
        }
    }
    drag(fluid) {

        var speed = this.velocity.magnitude();
        var dragMagnitude = fluid.c * speed * speed;

        var drag = this.velocity.get();
        drag.multiply(-1);
        drag.normalize();

        drag.multiply(dragMagnitude);

        this.applyForce(drag);
    }
}

var ballArr;
var liquid;
var liquidParam;


function reset() {
    ballArr = []
    for (i = 0; i < maxBalls; i++) {
        var moveInfo = {
            name_: 'particle ' + i,
            x_: random(0, width),
            y_: random(0, height),
            vx: random(-1, 1),
            vy: random(-1, 1),
            ax: random(-18, 18),
            ay: random(-18, 18),
            mass_: random(5, 20),
            topSpeed: 10,
        }
        ballArr.push(new Mover(moveInfo));
    }
}
function setup() {
    frameRate(50)
    createCanvas(innerWidth, innerHeight);
    ballArr = [];

    liquidParam = {
        x_: 0,
        y_: innerHeight / 2 + innerHeight / 3,
        w_: innerWidth,
        h_: innerHeight / 2,
        c_: .1,
    }
    pool = new Liquid(liquidParam);
    reset()
}
var epoch = 0;
var ts = 0;
function draw() {
    background(255);
    pool.display();
    for (var i = 0; i < ballArr.length; i++) {
        var wind = new Pvector(random(.05), 0);
        var gravity = new Pvector(0, 0.1 * (ballArr[i].mass));

        ballArr[i].applyForce(wind);
        ballArr[i].applyForce(gravity);

        ballArr[i].checkForWalls();
        ballArr[i].checkForCollision(ballArr);
        ballArr[i].display();
        ballArr[i].update();
        if (ballArr[i].isInside(pool)) {
            ballArr[i].drag(pool);
        }
        ballArr[i].limit();
    }
    ts++
    if (ts % 300 == 0) {
        epoch++
        ts = 0
        reset()
        console.log('new epoch')
    }
}