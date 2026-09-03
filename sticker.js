// sticker.js
// Sticker Maker на JavaScript (Node.js) с использованием @imgly/background-removal

const fs = require('fs');
const path = require('path');
const { removeBackground } = require('@imgly/background-removal');
const Jimp = require('jimp');

async function main() {
    const args = process.argv.slice(2);
    let input = null, output = null, color = null, threshold = 30, verbose = false;

    for (let i = 0; i < args.length; i++) {
        if (args[i] === '-o' || args[i] === '--output') {
            output = args[++i];
        } else if (args[i] === '-c' || args[i] === '--color') {
            color = args[++i];
        } else if (args[i] === '-t' || args[i] === '--threshold') {
            threshold = parseInt(args[++i]);
        } else if (args[i] === '-v' || args[i] === '--verbose') {
            verbose = true;
        } else if (args[i] === '-h' || args[i] === '--help') {
            console.log(`Использование: node sticker.js [опции] <входной_файл>
  -o, --output <файл>   Выходной файл (PNG)
  -c, --color <HEX>     Цвет фона (не используется)
  -t, --threshold <N>   Порог (не используется)
  -v, --verbose         Подробный вывод
  -h, --help            Справка`);
            process.exit(0);
        } else {
            input = args[i];
        }
    }

    if (!input) {
        console.error('Ошибка: не указан входной файл.');
        process.exit(1);
    }

    if (!fs.existsSync(input)) {
        console.error(`Ошибка: файл ${input} не найден`);
        process.exit(1);
    }

    output = output || input.replace(/\.[^.]+$/, '') + '_sticker.png';

    try {
        // Читаем изображение
        const image = await Jimp.read(input);
        // Конвертируем в base64 для библиотеки
        const base64 = await image.getBase64Async(Jimp.MIME_PNG);
        // Удаляем фон
        const resultBlob = await removeBackground(base64);
        // Преобразуем Blob в Buffer
        const buffer = Buffer.from(await resultBlob.arrayBuffer());
        fs.writeFileSync(output, buffer);
        if (verbose) {
            console.log(`Фон удалён. Результат сохранён в ${output}`);
        }
    } catch (err) {
        console.error(`Ошибка обработки: ${err.message}`);
        process.exit(1);
    }
}

main();
