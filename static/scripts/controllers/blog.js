'use strict';

angular.module('blogWebApp')
    .controller('BlogCtrl', ['$scope', '$http', '$window',function ($scope, $http , $window) {
        $scope.blogs = [];
        $scope.category = [];
        $scope.blogcategory = ""
        $scope.loading = "";
        $scope.blogeditkey = ""
        $scope.blogeditcat = ""
        $scope.headtitle = "Add Blog";
        $scope.working = false;
        
        var logError = function(data, status) {
            console.log('code '+status+': '+data);
            $scope.working = false;
            $scope.loading = "";
        };

        var refresh = function (){
            $scope.loading = "Loading Blogs ..."
            $http.get('/webservice/getallblog').
              success(function(data) { 
                $scope.loading = "";
                $scope.blogs = data.data; 

                if($scope.blogs == null && data.status == 0){
                  $scope.loading = "Blog not found.";
                }else if(data.status != 0){
                  $scope.loading = data.detail;
                }

              }).
              error(function(data,status){
                $scope.loading = "Cannot connect blog service , Error : " + status ;
                logError(data,status)
              });
              getcategory();
        };

        var getcategory = function (){
            $scope.blogeditcat = ""
            $scope.category=[]
            $http.get('/webservice/getallcategory').
              success(function(data) { 
                $scope.category = data.data; 
              }).
              error(logError);
        };

        $scope.canceledit = function(){
            $scope.blogTitle = "";
            $scope.blogContent = "";
            $scope.headtitle = "Add Blog"
            $scope.blogeditkey = ""

            getcategory();
        };

        $scope.saveblog = function (keyparam) {

            if(keyparam != ""){
              $http.post('/webservice/editblog/', {
                title: $scope.blogTitle,
                content: $scope.blogContent,
                category: $scope.blogcategory,
                key: keyparam
              }).
              error(logError).
              success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Edit Blog \n Status : " + data.status + "\n Detail : " + data.detail)
                  }else{
                    $scope.blogTitle = "";
                    $scope.blogContent = "";
                    $scope.headtitle = "Add Blog"
                    $scope.blogeditkey = ""
                    $scope.blogeditcat = ""
                    refresh();
                  }
              });
            }else{
              $http.post('/webservice/addblog/', {
                title: $scope.blogTitle,
                content: $scope.blogContent,
                btype: "blog",
                category: $scope.blogcategory,
              }).
              error(logError).
              success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Add Blog \n Status : " + data.status + "\n Detail : " + data.detail)
                  }else{
                    $scope.blogTitle = "";
                    $scope.blogContent = "";
                    refresh();
                  }
              });
            }
            
        };

        $scope.getblogedit = function (keyparam) {
              $scope.working = true;
              $http.get('/webservice/getblog',{
                params: { key: keyparam }
              }).
              success(function(data) { 
                $scope.working      = false;
                $scope.headtitle    = "Edit Blog"
                $scope.blogeditkey  = data.data[0].id
                $scope.blogeditcat  = data.data[0].value.category
                $scope.blogTitle    = data.data[0].value.title
                $scope.blogContent  = data.data[0].value.content
              }).
              error(logError);
              getcategory();
        };

        $scope.deleteblog = function (keyparam) {
            var deleteConfirm = $window.confirm('Are you absolutely sure you want to delete this blog?');

            if (deleteConfirm) {
                $http.get('/webservice/deleteblog',{
                  params: { key: keyparam }
                }).
                success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Delete Blog \n Status : " + data.status + "\n Detail : " + data.detail)
                  }else{
                    refresh()
                  }
                }).
                error(logError);
            }
        };

        refresh();
    }]);
