<template>
  <div class="home">
    <img alt="Vue logo" src="../assets/logo.png">
    <div style="width: 50%;margin: 0 auto">
      <el-form
          :model="address"
          ref="address"
          class="demo-ruleForm">
        <el-form-item
            label="Long Address"
        >
          <el-input
              v-model.number="address.url"
              autocomplete="off">

          </el-input>
          <el-input
              v-model.number="address.expiration_in_minutes"
              autocomplete="off">
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm('address')">提交</el-button>
          <el-button @click="resetForm('address')">重置</el-button>
          <el-button @click="searchForm('addressInfo')">查找</el-button>
        </el-form-item>
      </el-form>
    </div>

    <HelloWorld msg="Please Enter The Web Address"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue'
// import {longAdd} from "@/api/api";
import axios from "axios";

export default {
  name: 'HomeView',
  components: {
    HelloWorld
  },
  data() {
    return {
      address: {
        url: "",
        expiration_in_minutes:0,
      },
      shortLink:{
        short_link:""
      },
      addressInfo:[
    ]


    };
  },

  methods: {
    submitForm(formName) {
      this.$refs[formName].validate(async (valid) => {
        // console.log(valid)
        if (valid) {
          await this.longAdd();
        } else {
          console.log('error submit!!');
          return false;
        }
      });
      console.log()
    },
    searchForm(formName) {
      this.$refs[formName].validate(async (valid) => {
        // console.log(valid)
        if (valid) {
          await this.AddressInfo();
        } else {
          console.log('error submit!!');
          return false;
        }
      });
      console.log()
    },
    longAdd(){
      axios.post('http://localhost:9090/api/v1/long',this.address)
          .then((response)=>{
            console.log(response.data.short_link);
            this.shortLink= response.data
            console.log(this.shortLink);
          })
          .catch(function (error) {
            // handle error
            console.log(error);
          })
          .finally(function () {
            // always executed
          });
    },
    AddressInfo(){
      axios.get('http://localhost:9090/api/v1/info',this.shortLink)
          .then((response)=>{
            this.addressInfo = response.data
            console.log(this.addressInfo);

          })
          .catch(function (error) {
            // handle error
            console.log(error);
          })
          .finally(function () {
            // always executed
          });
    },
    resetForm(formName) {
      this.$refs[formName].resetFields();
    }
  }
}
</script>
