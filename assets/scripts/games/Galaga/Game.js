var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var frameCount = 0;
var fps, fpsInterval, startTime, now, then, elapsed;
fps = 60;
var moved = 0;
var ship, enemies, bullets;
var shipImage = new Image();
shipImage.src = "/assets/images/games/galaga/ship.png";
function startGame() {
  let shipY = canvas.height - 75;
  let shipX = canvas.width / 2 - 25;
  ship = new Rocket(shipX, shipY);
}

document.addEventListener('keydown', function(e) {
  if(moved === 0) {
    if(e.which === 37) {
      moved++;
      ship.moveLeft();
    } else if(e.which === 38) {
      moved++;
      ship.moveUp();
    } else if(e.which === 39) {
      moved++;
      ship.moveRight();
    } else if(e.which === 30) {
      moved++;
      ship.moveDown();
    }
  }
});
$(document).ready(function(){
  fpsInterval = 1000 / fps;
  then = Date.now();
  startTime = then;
  animate();
});
function animate() {
  requestAnimationFrame(animate);
  //Calculate elapsed time
  now = Date.now();
  elapsed = now - then;

  if (elapsed < fpsInterval) {
    //Dont draw
    return;
  }
  context.drawImage(shipImage, 0,0, 900, 900, ship.x, ship.y, 50, 50);
  moved = 0;


}
