class Rocket {
  constructor(xVal, yVal) {
    this.x = xVal;
    this.y = yVal;
    this.dx = 0;
    this.dy = 0;
    this.points = 0;
  }
  moveLeft() {
    this.x = this.x - 5;
  }
  moveRight() {
    this.x = this.x + 5;
  }
  moveUp() {
    this.y = this.y - 5;
  }
  moveDown() {
    this.y = this.y + 5;
  }
  shoot() {
    let blt = new Bullet(ship.x + 25, ship.y);
    return blt;
  }

}
