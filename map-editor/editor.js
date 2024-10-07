const canvas = document.querySelector("canvas");
const ctx = canvas.getContext("2d");

const exportButton = document.querySelector("button");

window.onload = () =>  {
    let map = new Image();
    map.src = "./diplomacy-map.png";

    map.onload = () => {
        canvas.width = map.width;
        canvas.height= map.height;

        ctx.drawImage(map, 0, 0, canvas.width, canvas.height);
    }
}

const circle = (x, y) => {
    const radius = 5;

    ctx.beginPath();
    ctx.arc(x, y, radius, 0, Math.PI*2);
    ctx.fillStyle = "pink";
    ctx.fill();
    ctx.strokeStyle = "red";
    ctx.stroke();
}

canvas.onclick = (event) => {
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    circle(x, y);
}

exportButton.onclick = () => {
    const url = canvas.toDataURL("image/png");
    const link = document.createElement("a");

    link.href = url;
    link.download = "diplomacy-map-edited.png";

    link.click();
}
