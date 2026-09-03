# sticker.rb
# Sticker Maker на Ruby (использует mini_magick)

require 'mini_magick'
require 'optparse'

options = {}
OptionParser.new do |opts|
  opts.banner = "Использование: ruby sticker.rb [опции] <входной_файл>"
  opts.on('-o', '--output FILE', 'Выходной файл (PNG)') { |v| options[:output] = v }
  opts.on('-c', '--color HEX', 'Цвет фона') { |v| options[:color] = v }
  opts.on('-t', '--threshold N', Integer, 'Порог (0-100)') { |v| options[:threshold] = v }
  opts.on('-v', '--verbose', 'Подробный вывод') { options[:verbose] = true }
  opts.on('-h', '--help', 'Справка') { puts opts; exit }
end.parse!

input = ARGV[0]
unless input
  puts "Ошибка: не указан входной файл."
  exit 1
end

unless File.exist?(input)
  puts "Ошибка: файл #{input} не найден"
  exit 1
end

output = options[:output] || input.sub(/\.[^.]+$/, '') + '_sticker.png'
threshold = options[:threshold] || 30
verbose = options[:verbose]

begin
  img = MiniMagick::Image.open(input)
  width = img.width
  height = img.height

  # Определяем цвет фона
  if options[:color]
    bg_color = MiniMagick::Image.read("xc:#{options[:color]}") { |c| c.format('png') }[0].get_pixels[0][0]
    bg_r, bg_g, bg_b = bg_color
  else
    # Автоопределение: средний цвет краёв
    pixels = []
    (0...width).each do |x|
      pixels << img.get_pixels[x][0]
      pixels << img.get_pixels[x][height-1]
    end
    (0...height).each do |y|
      pixels << img.get_pixels[0][y]
      pixels << img.get_pixels[width-1][y]
    end
    sr = sg = sb = 0
    pixels.each do |p|
      sr += p[0]; sg += p[1]; sb += p[2]
    end
    n = pixels.size
    bg_r = (sr / n).to_i
    bg_g = (sg / n).to_i
    bg_b = (sb / n).to_i
  end

  if verbose
    puts "Цвет фона: RGB(#{bg_r},#{bg_g},#{bg_b})"
  end

  max_dist = Math.sqrt(3 * 255 * 255)
  thresh = (threshold / 100.0) * max_dist

  # Создаём новое изображение с альфа-каналом
  # MiniMagick не поддерживает попиксельную установку альфа напрямую, поэтому используем `convert` через обработку.
  # Мы можем создать маску и применить её, но проще использовать подход с `-fuzz` и `-transparent`.
  # Однако это не даст точного контроля. Для демонстрации используем упрощённый подход с `-transparent-color`.
  # Но поскольку нам нужен точный порог, мы используем `-fuzz` и `-transparent`.
  # Альтернатива: использовать RMagick с попиксельным доступом, но это требует установки.

  # Вместо сложного попиксельного подхода, используем команду ImageMagick:
  # convert input.png -fuzz XX% -transparent "#RGB" output.png
  # Но это не позволяет указать цвет фона, отличный от одного цвета.
  # Для лучшего результата можно использовать несколько проходов.
  # Однако для демонстрации мы используем простой способ с `-transparent` и `-fuzz`.
  # Для более точного контроля в Ruby можно использовать RMagick.
  # Я предлагаю использовать упрощённый вариант: просто вызвать convert с -transparent.

  # Для корректной работы с RGB
  bg_hex = format('#%02x%02x%02x', bg_r, bg_g, bg_b)
  fuzz = threshold * 1.0
  # Создаём маску, но проще использовать -transparent
  # Делаем копию и применяем
  img2 = img.copy
  img2.combine_options do |c|
    c.fuzz("#{fuzz}%")
    c.transparent(bg_hex)
  end
  # Устанавливаем формат PNG и сохраняем
  img2.format('png')
  img2.write(output)
  img2.destroy

  if verbose
    puts "Фон удалён. Результат сохранён в #{output}"
  end
rescue => e
  puts "Ошибка: #{e.message}"
  exit 1
end
