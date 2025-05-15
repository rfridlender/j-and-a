# J&A Crew Management System
*J&A Crew Management System* is a web application designed to modernize and streamline the process of managing crew assignments, job details, and weekly payroll generation for contracting or service-based businesses.

The application addresses the need for a centralized system where crew leads can efficiently assign members to jobs, record crucial information such as addresses and site-specific notes, and ensure that this data seamlessly feeds into an automated payroll process. This aims to reduce manual data entry, minimize errors, and improve overall operational efficiency.

Technically, the system is built with a modern serverless architecture. The back end leverages Golang running on AWS Lambda, with data persisted in AWS DynamoDB utilizing a single-table design for optimized querying. Infrastructure is managed as code using Terraform, with distinct configurations for different environments (e.g., development, production) organized by directory for clarity and flexibility. Continuous integration and deployment (CI/CD) for both the front-end and back-end infrastructure are automated via GitHub Actions, deploying to AWS CloudFront.

## Technologies Used
### Front End
![TailwindCSS](https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=for-the-badge&logo=typescript&logoColor=white)
![Vue3](https://img.shields.io/badge/Vue3-002E3B?style=for-the-badge&logo=nuxtdotjs&logoColor=#00DC82)

![AWS Amplify](https://img.shields.io/badge/AWS-Amplify-ff9900?logo=awsamplify)

### Back End
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

![AWS API Gateway](https://img.shields.io/badge/AWS-API_Gateway-8C4fff?logo=amazonapigateway&logoColor=8C4fff)
![AWS CloudWatch](https://img.shields.io/badge/AWS-CloudWatch-dd344c?logo=amazoncloudwatch)
![AWS Cognito](https://img.shields.io/badge/AWS-Cognito-dd344c?logo=amazoncognito)
![AWS DynamoDB](https://img.shields.io/badge/AWS-DynamoDB-527fff?logo=amazondynamodb)
![AWS Lambda](https://img.shields.io/badge/AWS-Lambda-ff9900?logo=awslambda)

### DevOps
![GitHub Actions](https://img.shields.io/badge/GitHub-Actions-2088FF?logo=githubactions)
![Shell Script](https://img.shields.io/badge/shell_script-%23121011.svg?style=for-the-badge&logo=gnu-bash&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-844FBA?logo=terraform)

![AWS CloudFront](https://img.shields.io/badge/AWS-CloudFront-527fff)
![AWS S3](https://img.shields.io/badge/AWS-S3-569a31?logo=amazons3)

## Attributions
- UI Components inspired by and built with [shadcn-vue](https://www.shadcn-vue.com/)
- Icons from [Lucide Icons](https://lucide.dev/)
