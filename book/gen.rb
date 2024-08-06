#!/usr/bin/env ruby

require 'fileutils'
require 'optparse'

def gen
  `protoc --go_out=./ *.proto;mv iting.com/book/* .; rm -rf iting.com/`
  puts "生成成功"
  # deploy
end
gen
# def deploy
#   # 获取当前目录中所有的 .go 文件
#   go_files = Dir.glob("#{File.dirname(__FILE__)}/*.{go,proto}")

#   # 指定目标目录
#   target_directory = File.expand_path("~/dev/go/pocketbase/book")

#   # 确保目标目录存在
#   unless Dir.exist?(target_directory)
#     # target_directory parent directory should exist
#     if !Dir.exist?(File.dirname(target_directory))
#       puts "目标目录不存在: #{target_directory}"
#       exit
#     end
#     FileUtils.mkdir_p(target_directory)
#   end

#   # 遍历 .go 文件并创建硬链接
#   go_files.each do |file|
#     target_file = File.join(target_directory, File.basename(file))
#     begin
#       # 如果目标文件存在，先删除它
#       FileUtils.rm_f(target_file)
#       # 创建硬链接
#       FileUtils.ln(file, target_file)
#       puts "硬链接创建成功: #{file} -> #{target_file}"
#     rescue => e
#       puts "创建硬链接失败: #{file} -> #{target_file}, 错误: #{e.message}"
#     end
#   end

# end


# def run
#   # 获取命令行参数
#   command = ARGV[0]

#   # 根据命令行参数路由到适当的方法
#   case command
#   when "gen"
#     gen
#   when "deploy"
#     deploy
#   else
#     puts "命令无效: #{command}"
#     puts "请指定一个有效的命令: gen 或 deploy"
#     puts "用法: ruby your_script.rb [gen|deploy]"
#   end
# end

# run
