class Bullet {
  constructor(xVal, yVal) {
    this.x = xVal;
    this.y = yVal;
    this.dy = -15;
  }
  move() {
    this.y += this.dy;
  }
}
