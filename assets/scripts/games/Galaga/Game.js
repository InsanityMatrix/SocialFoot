var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var frameCount = 0;
var fps, fpsInterval, startTime, now, then, elapsed;
fps = 60;
var stars = [];
var gameState;
var moved = 0;
var round;
var keysPressed = {
  forward: false,
  left: false,
  right: false,
  backward: false,
  shoot: false
}
var fire = 0;
var fireTimer = 0;
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
  fire = 0;
  fireTimer = 0;
  for(var i = 0; i < 75; i++) {
    let xVal = Math.floor(Math.random() * canvas.width);
    let yVal = Math.floor(Math.random() * canvas.height);
    stars[i] = new Star(xVal,yVal,canvas.height);
  }
  round = 1;
  newWave(round);
}

document.addEventListener('keydown', function(e) {
    if(e.which === 37) {
      e.preventDefault();
      keysPressed.left = true;
    } else if(e.which === 38) {
      e.preventDefault();
      keysPressed.forward = true;
    } else if(e.which === 39) {
      e.preventDefault();
      keysPressed.right = true;
    } else if(e.which === 40) {
      e.preventDefault();
      keysPressed.backward = true;
    }

  if(e.which === 70) {
    if(gameState === 2 || gameState === 1) {
      gameState = 0;
    } else if (gameState === 3){
      startGame();
      gameState = 0;
    } else {
      gameState = 1;
    }
  } else if(e.which === 32) {
    e.preventDefault();
    keysPressed.shoot = true;

  }

});
document.addEventListener('keyup', function(e) {
    if(e.which === 37) {
      keysPressed.left = false;
    } else if(e.which === 38) {
      keysPressed.forward = false;
    } else if(e.which === 39) {
      keysPressed.right = false;
    } else if(e.which === 40) {
      keysPressed.backward = false;
    }

   if(e.which === 32) {
    keysPressed.shoot = false;
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
  if(fireTimer++ === 4) {
    fireTimer = 0;
    fire = 0;
  }
  if(keysPressed.forward) {
    if(ship.y > 0) {
      ship.moveUp();
    }
  } else if(keysPressed.backward) {
    if(!(ship.y + 50 >= canvas.height)) {
      ship.moveDown();
    }
  }
  if(keysPressed.left) {
    if(ship.x > 0) {
      ship.moveLeft();
    }
  } else if(keysPressed.right) {
    if(!(ship.x + 50 >= canvas.width)) {
      ship.moveRight();
    }
  }
  if(keysPressed.shoot) {
    if(fire === 0) {
      bullets[bullets.length] = ship.shoot();
      bullets[bullets.length - 1].move();
      fire++;
    }
  }
  context.clearRect(0,0,canvas.width,canvas.height);
  context.fillStyle = "black";
  context.fillRect(0,0, canvas.width, canvas.height);
  context.fillStyle = "white";
  for(var i = 0; i < stars.length; i++) {
    context.fillRect(stars[i].x, stars[i].y, 2,2);
    stars[i].move();
  }

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
  } else if(gameState === 3) {
    context.font = "35px Comic Sans MS";
    context.fillStyle = "white";
    context.textAlign = "center";
    context.fillText("You Lost!", canvas.width / 2, 50);
    context.fillText("Wave: " + round, canvas.width / 2, 90);
    context.fillText("(Press F to play again)", canvas.width /2, canvas.height -10);
    return;
  }

  //Draw ship
  context.drawImage(shipImage, 0,0, 900, 900, ship.x, ship.y, 50, 50);
  moved = 0;
  //Draw Enemies
  for(var i = 0; i < enemies.length; i++) {
    let enemy = enemies[i];
    if(!enemy.dead) {
      if(enemy.y >= canvas.height - 5) {
        gameState = 3;
        return;
      }
      if(enemy.level === 1) {
        context.drawImage(enemyL1, 0, 0, 305, 240, enemy.x, enemy.y, 20, 20);
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
      context.drawImage(bulletIMG, 0, 0, 407,512, bullets[i].x - 4, bullets[i].y,8,12);
    }
  }
  //Check if enemies killed
  for(var enemyIndex = 0; enemyIndex < enemies.length; enemyIndex++) {
    for(var bulletIndex = 0; bulletIndex < bullets.length; bulletIndex++) {
      let e = enemies[enemyIndex];
      let bullet = bullets[bulletIndex];

      if(bullet.x - e.x >=0 &&
         bullet.x - e.x <= 25 &&
         bullet.y - e.y <= 25 &&
         bullet.y - e.y >= 0 &&
         !e.dead) {
         enemies[enemyIndex].hit();
         if(bulletIndex === 0) {
           bullets.splice(0,1);
         } else {
           bullets.splice(bulletIndex, bulletIndex);
         }
      }
    }
  }
  for(var i = 0; i < enemies.length; i++) {
    let enemy = enemies[i];
    if(!enemy.dead) {
      if(enemy.y - ship.y >= 0 &&
         enemy.y - ship.y < 50 &&
         enemy.x - ship.x >= 0 &&
         enemy.x - ship.x < 50) {
           gameState = 3;
           return;
      }
    }
  }
    //Check if all enemies are dead if so new wave
    for (var i = 0; i < enemies.length; i++) {
      if(!enemies[i].dead) {
        break;
      }
      if(i === enemies.length - 1) {
        // All enemies are dead;
        enemies = [];
        setTimeout(newWave, 500, ++round);
      }
    }


  //Move all bullets
  for(var i = 0; i < bullets.length; i++) {
    bullets[i].move();
  }
  //Move all enemies
  for(var i = 0; i < enemies.length; i++) {
    enemies[i].x += enemies[i].dx;
  }
}
function newWave(num) {
  if(num === 1) {
    let newEnemies = [];
    for(var i = 0; i < 20; i++) {
      let e = new Enemy(1, i,canvas.width, num);
      newEnemies[i] = e;
    }
    enemies = newEnemies;
  } else {
    let newEnemies = [];
    for(var i = 0; i < 25; i++) {
      let e = new Enemy(1, i,canvas.width, num);
      newEnemies[i] = e;
    }
    enemies = newEnemies;
  }
}
