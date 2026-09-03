// sticker.ts
// Sticker Maker на TypeScript (использует @imgly/background-removal)

import * as fs from 'fs';
import { removeBackground } from '@imgly/background-removal';
import Jimp from 'jimp';

async function main(): Promise<void> {
    const args = process.argv.slice(2);
    let input: string | null = null;
    let output: string | null = null;
    let color: string | null = null;
    let threshold: number = 30;
    let verbose: boolean = false;

    for (let i = 0; i < args.length; i++) {
        switch (args[i]) {
            case '-o':
            case '--output':
                output = args[++i];
                break;
            case '-c':
            case '--color':
                color = args[++i];
                break;
            case '-t':
            case '--threshold':
                threshold = parseInt(args[++i]);
                break;
            case '-v':
            case '--verbose':
                verbose = true;
                break;
            case '-h':
            case '--help':
                console.log(`Использование: ts-node sticker.ts [опции] <входной_файл>
  -o, --output <файл>   Выходной файл (PNG)
  -c, --color <HEX>     Цвет фона (не используется)
  -t, --threshold <N>   Порог (не используется)
  -v, --verbose         Подробный вывод
  -h, --help            Справка`);
                process.exit(0);
            default:
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
        const image = await Jimp.read(input);
        const base64 = await image.getBase64Async(Jimp.MIME_PNG);
        const resultBlob = await removeBackground(base64);
        const buffer = Buffer.from(await resultBlob.arrayBuffer());
        fs.writeFileSync(output, buffer);
        if (verbose) {
            console.log(`Фон удалён. Результат сохранён в ${output}`);
        }
    } catch (err: any) {
        console.error(`Ошибка обработки: ${err.message}`);
        process.exit(1);
    }
}

main();
