var gameboard;
var grid;
var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var score = 0;
var gameState = 0;

// gameState
// 0 = Playing
// 1 = Lost

function startGame() {
  gameboard = [
    [0,0,0,0],
    [0,0,0,0],
    [0,0,0,0],
    [0,0,0,0]
  ];
  let cell = chooseEmptyCell();
  gameboard[cell[0]][cell[1]] = 2;
  cell = chooseEmptyCell();
  gameboard[cell[0]][cell[1]] = 2;
}

function chooseEmptyCell() {
  let tries = 0;
  while(tries < 125) {
    let rRow = Math.floor(Math.random() * 4);
    let rCol = Math.floor(Math.random() * 4);

    let cell = gameboard[rRow][rCol];
    if(cell === 0) {
      return [rRow, rCol];
    }
    tries++;
  }
  return [-1, -1];
}
function moveLeft() {
  //Start combine from left to right
  for(var i = 0; i < 4; i++) {
    for(var j = 0; j < 3; j++) {
      if(gameboard[i][j] === gameboard[i][j + 1]) {
        gameboard[i][j] *= 2;
        gameboard[i][j + 1] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var i = 0; i < 4; i++) {
    for(var j = 3; j > 0; j--) {
      if(gameboard[i][j - 1] === 0 && gameboard[i][j] != 0) {
        //Shift this cell left
        gameboard[i][j - 1] = gameboard[i][j];
        gameboard[i][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveRight() {
  for(var i = 0; i < 4; i++) {
    for(var j = 3; j > 0; j--) {
      if(gameboard[i][j] === gameboard[i][j - 1]) {
        gameboard[i][j] *= 2;
        gameboard[i][j - 1] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var i = 0; i < 4; i++) {
    for(var j = 0; j < 3; j++) {
      if(gameboard[i][j + 1] === 0 && gameboard[i][j] != 0) {
        //Shift this cell right
        gameboard[i][j + 1] = gameboard[i][j];
        gameboard[i][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveUp() {
  for(var j = 0; j < 4; j++) {
    for(var i = 0; i < 3; i++) {
      if(gameboard[i][j] === gameboard[i + 1][j]) {
        gameboard[i][j] *= 2;
        gameboard[i+1][j] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var j = 0; j < 4; j++) {
    for(var i = 3; i > 0; i--) {
      if(gameboard[i-1][j] === 0 & gameboard[i][j] != 0) {
        gameboard[i-1][j] = gameboard[i][j];
        gameboard[i][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveDown() {
  for(var j = 0; j < 4; j++) {
    for(var i = 3; i > 0; i--) {
      if(gameboard[i][j] === gameboard[i - 1][j]) {
        gameboard[i][j] *= 2;
        gameboard[i-1][j] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var j = 0; j < 4; j++) {
    for(var i = 0; i < 3; i++) {
      if(gameboard[i+1][j] === 0 & gameboard[i][j] != 0) {
        gameboard[i+1][j] = gameboard[i][j];
        gameboard[i][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function spawnNewCell() {
  let value = 2;
  if(Math.random() > 0.7) {
    value = 4;
  }
  cell = chooseEmptyCell();
  if (cell[0] === -1) {
    gameSate = 1;
    return;
  }
  gameboard[cell[0]][cell[1]] = value;
}
document.addEventListener('keydown', function(e) {
  //LEFT
  if(e.which === 37) {
    moveLeft();
  }
  //RIGHT
  if(e.which === 39) {
    moveRight();
  }
  //UP
  if(e.which === 38) {
    moveUp();
  }
  //DOWN
  if(e.which === 40) {
    moveDown();
  }
});
//Precondition width and height = 500
function loop() {
  requestAnimationFrame(loop);
  //60 FPS 2048 Babyyy
  if(gameState === 0) {
    context.clearRect(0,0,canvas.width, canvas.height);
    //Borders = 20
    context.fillStyle = "#6e4f31";
    //Vertical lines
    context.fillRect(0,0,20,canvas.height);
    context.fillRect(120,0,20,canvas.height);
    context.fillRect(240, 0, 20, canvas.height);
    context.fillRect(360,0,20,canvas.height);
    context.fillRect(480,0,20,canvas.height);
    //Horizontal Lines
    context.fillRect(0,0,canvas.width, 20);
    context.fillRect(0,120,canvas.width, 20);
    context.fillRect(0,240,canvas.width, 20);
    context.fillRect(0,360,canvas.width, 20);
    context.fillRect(0,480,canvas.width,20);

    return;
  }
}
startGame();
requestAnimationFrame(loop);
