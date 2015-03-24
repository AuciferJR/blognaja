'use strict';

angular.module('blogWebApp')
    .controller('CategoryCtrl', ['$scope', '$http', '$window',function ($scope, $http , $window) {
        $scope.blogs = [];
        $scope.loading = "";
        $scope.blogeditkey = ""
        $scope.headtitle = "Add Category";
        $scope.working = false;
        
        var logError = function(data, status) {
            console.log('code '+status+': '+data);
            $scope.working = false;
            $scope.loading = "";
        };

        var refresh = function (){
            $scope.loading = "Loading Category ..."
            $http.get('/webservice/getallcategory').
              success(function(data) { 
                $scope.loading = "";
                $scope.blogs = data.data; 

                if($scope.blogs == null && data.status == 0){
                  $scope.loading = "Category not found.";
                }else if(data.status != 0){
                  $scope.loading = data.detail;
                }

              }).
              error(function(data,status){
                $scope.loading = "Cannot connect blog service , Error : " + status ;
                logError(data,status)
              });
        };

        $scope.canceledit = function(){
            $scope.blogTitle = "";
            $scope.blogContent = "";
            $scope.headtitle = "Add Category"
            $scope.blogeditkey = ""
        };

        $scope.saveblog = function (keyparam) {

            if(keyparam != ""){
              $http.post('/webservice/editblog/', {
                title: $scope.blogTitle,
                content: $scope.blogContent,
                key: keyparam
              }).
              error(logError).
              success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Edit Category \n Status : " + data.status + "\n Detail : " + data.detail)
                  }else{
                    $scope.blogTitle = "";
                    $scope.blogContent = "";
                    $scope.headtitle = "Add Category"
                    $scope.blogeditkey = ""
                    refresh();
                  }
              });
            }else{
              $http.post('/webservice/addblog/', {
                title: $scope.blogTitle,
                content: $scope.blogContent,
                btype: "category"
              }).
              error(logError).
              success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Add Category \n Status : " + data.status + "\n Detail : " + data.detail)
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
                $scope.working = false;
                $scope.headtitle = "Edit Category"
                $scope.blogeditkey = data.data[0].id
                $scope.blogTitle = data.data[0].value.title
                $scope.blogContent = data.data[0].value.content
              }).
              error(logError);
        };

        $scope.deleteblog = function (keyparam) {
            var deleteConfirm = $window.confirm('Are you absolutely sure you want to delete this Category?');

            if (deleteConfirm) {
                $http.get('/webservice/deleteblog',{
                  params: { key: keyparam }
                }).
                success(function(data) {
                  if(data.status != 0){
                      $window.alert("Cannot Delete Category \n Status : " + data.status + "\n Detail : " + data.detail)
                  }else{
                    refresh()
                  }
                }).
                error(logError);
            }
        };

        refresh();
    }]);
