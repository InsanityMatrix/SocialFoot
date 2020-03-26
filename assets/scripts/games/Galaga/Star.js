class Star {
  constructor(xVal, yVal, height) {
    this.x = xVal;
    this.y = yVal;
    this.height = height;
  }
  move() {
    this.y = (this.y + 10) % this.height;
  }
}
