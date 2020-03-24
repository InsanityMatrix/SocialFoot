class Enemy {
  constructor(lvl, ct, width, wave) {
    this.level = lvl;
    this.hp = lvl * 10;
    this.dead = false;
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
}
