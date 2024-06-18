#!/bin/bash
set -euxo pipefail
 
#-------------------------------------------------------------------------------------------
# おまじない
#-e: エラーが発生したときにスクリプトを中断する。途中でエラーにしたくない場合はset +eで一時的に解除するか||で繋げる
#-u: 未定義変数をエラーにする
#-x: 実行したコマンドを出力する
#-o pipefail: パイプで結合したコマンドの途中でエラーが発生したときもスクリプトを中断する
#---------------------------------------------------------------------------------------------
 
#VTLのデータセット名と最大世代数情報を取得
python ./put_VTLInfo.py
 
 
#VTLカタログ情報からデータセットとVOLSER情報を抽出
python ./put_VTLInfo_DATASET_VOL.py
 
#VOLSER名でソート
sort -t ',' -k 2 VTLInfo_daset_volser.csv > sorted_VTLInfo_daset_volser.csv
 
#VOLSER単位でデータセットを整理してcsv化
python ./dup_VOLSER_DATASET.py