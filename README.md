# wget_batch_rename
多线程批量使用wget下载文件，并自动重命名

使用方法：
  
   * go run wget_batch.go -ext jpg -i inputfile -prefix="eweb"
   
   * inputfile: 批量下载资源，text文件，每行一个url
   
   * 下载文件重命名规则: [prefix]+[index].[ext]
