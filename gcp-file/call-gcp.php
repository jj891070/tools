<?php
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);
//判斷是否有該檔案
$filename="yes.mp3";
$str=file_get_contents($filename);
// if(file_exists($filename)){

//     $file = fopen($filename, "r");

//     if($file != NULL){

//         //當檔案未執行到最後一筆，迴圈繼續執行(fgets一次抓一行)

//         while (!feof($file)) {

//             $str .= fgets($file);

//         }

//         fclose($file);

//     }

// }
$encode=base64_encode($str);
// exec("GOOGLE_APPLICATION_CREDENTIALS='/Users/jay_liao/go/src/jay-test/in-sports-test-storage.json' GOOGLE_CLOUD_PROJECT='in-sports-prod' ./jay-test -o=jay-test-fat-log:zoular/1/2/jay.json  write '${encode}'", $output, $return_var);
   
$curl = curl_init();
$write=[
    "method"=>"GcpFile.Upload",
    // "service"=>"gcp-file",
    "params"=>[
        "remote_bucket"=>"jay-test-fat-log",
        "data"=>[
            "annie_game/1/123.mp3"=>$encode,
            "annie_game2/1/123.mp3"=>$encode,
            "annie_game3/1/123.mp3"=>$encode,
            "annie_game4/1/123.mp3"=>$encode,
            "annie_game5/1/123.mp3"=>$encode,
            "annie_game6/1/123.mp3"=>$encode,
            "annie_game7/1/123.mp3"=>$encode,
            "annie_game8/1/123.mp3"=>$encode,
            "annie_game9/1/123.mp3"=>$encode,
            
        ]
    ]
        ];

$read=[
    "method"=>"GcpFile.Read",
    "service"=>"gcp-file",
    "params"=>[
        "remote_bucket"=>"jay-test-fat-log",
        "remote_path"=>[
            "jay.json",
            "jay/1/2/jay1.json",

        ]
    ]
        ];

$delete=[
    "method"=>"GcpFile.Delete",
    // "service"=>"gcp-file",
    "params"=>[
        "remote_bucket"=>"jay-test-fat-log",
        "remote_path"=>[
            "沙巴賽事.txt"
        ]
    ]
        ];
curl_setopt_array($curl, array(
  CURLOPT_PORT => "8000",
  CURLOPT_URL => "http://127.0.0.1:8000",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "POST",
  CURLOPT_POSTFIELDS => json_encode($write),
  CURLOPT_HTTPHEADER => array(
    "cache-control: no-cache",
    "postman-token: 100042d9-0934-757f-edf2-250dfb4b3d8e"
  ),
));

$response = curl_exec($curl);
$err = curl_error($curl);

curl_close($curl);

if ($err) {
  echo "cURL Error #:" . $err;
} else {
  echo $response;
}

