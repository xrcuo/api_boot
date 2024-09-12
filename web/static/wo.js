
$(document).ready(function () {
    setInterval(updateSystemInfo, 1000);
    setInterval(updateNetworkSpeed, 1000);
    setInterval(updateProcessInfo, 1000);
});

function updateSystemInfo() {
    $.getJSON("/system_info", function (data) {
        // 更新其他系统信息
        $("#cpu-percent").text(data.cpuPercent.toFixed(2));
        $("#memory-used").text(data.memoryUsed.toFixed(2));

        // 更新硬盘信息表格
        var diskInfoTbody = $("#disk-info tbody");
        diskInfoTbody.empty();
        $.each(data.diskInfos, function (index, disk) {
            diskInfoTbody.append("<tr><td>" + disk.mountPoint + "</td><td>" + disk.total + "</td><td>" + disk.used + "</td><td>" + disk.usedPercent.toFixed(2) + "%</td></tr>");
        });

        // 更新 CPU 仪表盘
        updateCPUGauge(data.cpuPercent);

        // 更新内存仪表盘
        updateMemoryGauge(data.memoryUsed);
    });
}

// 更新 CPU 仪表盘
function updateCPUGauge(cpuUsage) {
    const cpuGaugeValue = document.getElementById('cpuGaugeValue');
    cpuGaugeValue.textContent = `${cpuUsage.toFixed(1)}%`;
}

// 更新内存仪表盘
function updateMemoryGauge(memoryUsage) {
    const memoryGaugeValue = document.getElementById('memoryGaugeValue');
    memoryGaugeValue.textContent = `${memoryUsage.toFixed(1)}%`;
}

function updateNetworkSpeed() {
    // 获取速度单位
    var speedUnit = "{{.SpeedUnit}}";

    // 发送 AJAX 请求获取网络速度信息
    $.getJSON("/speed", function (data) {
        // 清空表格内容
        $("#network-interfaces tbody").empty();

        // 遍历网络接口数据
        for (var key in data) {
            if (data.hasOwnProperty(key)) {
                var speed = data[key];

                // 格式化速度值
                var downSpeed = formatSpeed(speed.BytesRecv, speedUnit);
                var upSpeed = formatSpeed(speed.BytesSent, speedUnit);

                // 添加新行到表格中
                $("#network-interfaces tbody").append("<tr><td>" + speed.Name + "</td><td>" + downSpeed + "</td><td>" + upSpeed + "</td></tr>");
            }
        }
    });
}

function updateProcessInfo() {
    $.getJSON("/processes", function (data) {
        // 更新进程信息表格
        var processInfoTbody = $("#process-info tbody");
        processInfoTbody.empty();
        $.each(data, function (index, process) {
            processInfoTbody.append("<tr><td>" + process.pid + "</td><td>" + process.name + "</td><td>" + process.cpuPercent.toFixed(2) + "</td><td>" + process.memoryUsed.toFixed(2) + "</td></tr>");
        });
    });
}

// 格式化速度值
function formatSpeed(bytes, unit) {
    // 如果字节数为 0，则直接返回 0
    if (bytes === 0) {
        return "0 " + unit;
    }

    var speed;
    if (unit === "mbps") {
        speed = bytes * 8 / 1000 / 1000; // Mbps
    } else {
        speed = bytes * 8 / 1000; // Kbps
    }
    return speed.toFixed(2) + " " + unit;
}
