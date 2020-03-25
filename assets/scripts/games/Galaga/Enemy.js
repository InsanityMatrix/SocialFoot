class Enemy {
  constructor(lvl, ct, width, wave) {
    this.level = lvl;
    if(lvl === 2) {
      this.inDive = false;
    }
    this.hp = lvl * 10;
    this.dead = false;
    this.width = width;
    let rowCount = ct % 5
    this.dx = wave * 2 % 10 + 1;
    if(wave < 10) {
      this.dx = wave;
    } else if(wave / 10 === 0) {
      this.dx = 12;
    }
    this.x = (5 + (25* rowCount));
    this.y = 5 + (25 * Math.floor(ct / 5));
  }
  switchDirections() {
    this.dx *= -1;
  }
  hit() {
    this.hp -= 10;
    if(this.hp <= 0) {
      this.dead = true;
    }
  }
  dive() {
    this.targetx = 0.5 * this.width;
    this.targety = this.width - 10;
    this.inDive = true;
    this.reflectedx = this.width - this.x;
    this.reflectedy = this.y;

    //Calculate Parabola
    let e1A = Math.pow(this.x, 2);
    let e1B = this.x;
    let e1Eq = this.y;

    let e2A = Math.pow(this.targetx, 2);
    let e2B = this.targetx;
    let e2Eq = this.targety;

    let e3A = Math.pow(this.reflectedx, 2);
    let e3B = this.reflectedx;
    let e3Eq = this.reflectedy;

    
  }

}
