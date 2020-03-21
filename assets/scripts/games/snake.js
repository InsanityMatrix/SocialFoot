var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var paused = false;
var grid = 16;
var count = 0;
var userid = 0;
var moves = 0;

var snake = {
  x: 160, y: 160, dx: grid, dy: 0,cells: [], maxCells: 4
}
var apple = {
  x: 320, y: 320
};

function setUserID(uid) {
  userid = uid;
}
function getRandomInt(min, max) {
  return Math.floor(Math.random() * (max - min)) + min;
}
function scorePage() {
  context.clearRect(0,0,canvas.width,canvas.height);
  context.font = "20px Comic Sans MS";
  context.fillStyle = "black";
  context.textAlign = "center";
  context.fillText("Scores", canvas.width/2, 30);

  $.ajax({
    url: '/games/snake/scores',
    success: populateScores
  });

  context.font = "15px Comic Sans MS";
  context.fillText("(Press Space to Continue)", canvas.width/2, canvas.height - 10);
}
function populateScores(data) {
  context.font = "20px Comic Sans MS";
  context.fillStyle = "black";
  context.textAlign = "center";

  context.fillText("1. " + data[0].Username, canvas.width/4 * 1.5, 60)
  context.fillText("" + data[0].Score, canvas.width/4 * 3, 60)

  context.fillText("2. " + data[1].Username, canvas.width/4 * 1.5, 90)
  context.fillText("" + data[1].Score, canvas.width/4 * 3, 90)

  context.fillText("3. " + data[2].Username, canvas.width/4 * 1.5, 120)
  context.fillText("" + data[2].Score, canvas.width/4 * 3, 120)

  context.fillText("4. " + data[3].Username, canvas.width/4 * 1.5, 150)
  context.fillText("" + data[3].Score, canvas.width/4 * 3, 150)

  context.fillText("5. " + data[4].Username, canvas.width/4 * 1.5, 180)
  context.fillText("" + data[4].Score, canvas.width/4 * 3, 180)
}
function loop() {
  requestAnimationFrame(loop);

  if(paused) {

    return;
  }
  //Make game 15 fps
  if (++count < 6) {
    return;
  }
  moves = 0;
  count = 0;
  context.clearRect(0,0,canvas.width,canvas.height);
  snake.x += snake.dx;
  snake.y += snake.dy;

  if (snake.x < 0) {
    restartGame();
    return;
  } else if (snake.x > canvas.width) {
    restartGame();
    return;
  }
  if (snake.y < 0) {
    restartGame();
    return;
  } else if (snake.y > canvas.height) {
    restartGame();
    return;
  }
  snake.cells.unshift({x: snake.x, y: snake.y});

  if (snake.cells.length > snake.maxCells) {
    snake.cells.pop();
  }

  context.fillStyle = 'red';
  context.fillRect(apple.x, apple.y, grid-1, grid-1);

  context.fillStyle = 'green';
  snake.cells.forEach(function(cell, index) {
    context.fillRect(cell.x, cell.y, grid-1, grid-1);
    if (cell.x === apple.x && cell.y === apple.y) {
      snake.maxCells++;

      apple.x = getRandomInt(0, 25) * grid;
      apple.y = getRandomInt(0, 25) * grid;
    }
    for (var i = index + 1; i < snake.cells.length; i++) {
      if (cell.x === snake.cells[i].x && cell.y === snake.cells[i].y) {
        //Reset game - todo change this
        restartGame();
        return;
      }
    }
  });
}

document.addEventListener('keydown', function(e) {
  // left arrow key
  if (e.which === 37 && snake.dx === 0) {
    if(moves > 0) {
      return;
    }
    snake.dx = -grid;
    snake.dy = 0;
    moves++;
  }
  // up arrow key
  else if (e.which === 38 && snake.dy === 0) {
    if(moves > 0) {
      return;
    }
    snake.dy = -grid;
    snake.dx = 0;
    moves++;
  }
  // right arrow key
  else if (e.which === 39 && snake.dx === 0) {
    if(moves > 0) {
      return;
    }
    snake.dx = grid;
    snake.dy = 0;
    moves++;
  }
  // down arrow key
  else if (e.which === 40 && snake.dy === 0) {
    if(moves > 0) {
      return;
    }
    snake.dy = grid;
    snake.dx = 0;
    moves++;
  }
  else if (e.which === 32) {
    if(paused) {
      paused = false;
    } else {
      paused = true;
    }
  }
});

function restartGame() {
  $.ajax({
    url: '/games/snake/update',
    method: 'POST',
    data: {
      "userid": userid,
      "score": snake.maxCells
    },
    success: scorePage
  });

  snake.x = 160;
  snake.y = 160;
  snake.cells = [];
  snake.maxCells = 4;
  snake.dx = grid;
  snake.dy = 0;

  apple.x = getRandomInt(0, 25) * grid;
  apple.y = getRandomInt(0, 25) * grid;
  paused = true;
}

requestAnimationFrame(loop);
