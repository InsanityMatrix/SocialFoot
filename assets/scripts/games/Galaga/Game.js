var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var frameCount = 0;
var fps, fpsInterval, startTime, now, then, elapsed;
fps = 60;
var gameState;
var moved = 0;
var ship, enemies, bullets;
bullets = [];
var shipImage = new Image();
shipImage.src = "/assets/images/games/galaga/ship.png";
var enemyL1 = new Image();
enemyL1.src = "/assets/images/games/galaga/lvl1Enemy.png";
var bulletIMG = new Image();
bulletIMG.src = "/assets/images/games/galaga/bullet.png";
function startGame() {
  let shipY = canvas.height - 75;
  let shipX = canvas.width / 2 - 25;
  ship = new Rocket(shipX, shipY);
  gameState = 2;
  newWave(1);
}

document.addEventListener('keydown', function(e) {
  if(moved === 0) {
    if(e.which === 37) {
      if(ship.x <= 0) {
        return;
      }
      moved++;
      ship.moveLeft();
    } else if(e.which === 38) {
      if(ship.y <= 0) {
        return;
      }
      moved++;
      ship.moveUp();
    } else if(e.which === 39) {
      if(ship.x + 50 >= canvas.width) {
        return;
      }
      moved++;
      ship.moveRight();
    } else if(e.which === 40) {
      if(ship.y + 50 >= canvas.height) {
        return;
      }
      moved++;
      ship.moveDown();
    }
  }
  if(e.which === 70) {
    if(gameState === 2 || gameState === 1) {
      gameState = 0;
    } else {
      gameState = 1;
    }
  } else if(e.which === 32) {
    bullets[bullets.length] = ship.shoot().move();
  }

});
$(document).ready(function(){
  fpsInterval = 1000 / fps;
  then = Date.now();
  startTime = then;
  startGame();
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
  context.clearRect(0,0,canvas.width,canvas.height);
  context.fillStyle = "black";
  context.fillRect(0,0, canvas.width, canvas.height);
  //Draw ship
  context.drawImage(shipImage, 0,0, 900, 900, ship.x, ship.y, 50, 50);
  moved = 0;

  if(gameState === 2) {
    context.font = "35px Comic Sans MS";
    context.fillStyle = "white";
    context.textAlign = "center";
    context.fillText("Press F to start the game.", canvas.width / 2, 50);
    return;
  } else if(gameState === 1) {
    context.font = "35px Comic Sans MS";
    context.fillStyle = "white";
    context.textAlign = "center";
    context.fillText("Press F to unpause", canvas.width / 2, 50);
    return;
  }
  //Draw Enemies
  for(var i = 0; i < enemies.length; i++) {
    let enemy = enemies[i];
    if(!enemy.dead) {
      if(enemy.level === 1) {
        context.drawImage(enemyL1, 0, 0, 305, 240, enemy.x, enemy.y, 20, 20);
        enemies[i].x += enemies[i].dx;
      }
    }
  }
  if(enemies[4].dx > 0) {
    if(Math.abs(canvas.width - (enemies[4].x + 25))  <= 5) {
      for(var i = 0; i < enemies.length; i++) {
        enemies[i].y += 25;
        enemies[i].switchDirections();
      }
    }
  } else {
    if(enemies[0].x <= 5) {
      for(var i = 0; i < enemies.length; i++) {
        enemies[i].y += 25;
        enemies[i].switchDirections();
      }
    }
  }
  //Draw Bullets
  for(var i = 0; i < bullets.length; i++) {
    if(bullets[i].y <= 0) {
      if(i === 0) {
        bullets.splice(0,1);
      } else {
        bullets.splice(i,i);
      }
    } else {
      context.drawImage(bulletIMG, 0, 0, 407,512, bullets[i].x - 3, bullets[i].y,6,6);
      bullets[i].move();
    }
  }
}
function newWave(num) {
  if(num === 1) {
    let newEnemies = [];
    for(var i = 0; i < 20; i++) {
      let e = new Enemy(1, i,canvas.width);
      newEnemies[i] = e;
    }
    enemies = newEnemies;
  }
}
