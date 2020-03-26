class Star {
  constructor(xVal, yVal, height) {
    this.x = xVal;
    this.y = yVal;
    this.height = height;
  }
  move() {
    this.y = this.height % (this.y + 10);
  }
}
