# open-seattle-example

## Utilities:

- We use [Just](https://github.com/casey/just) for quick access to some scripts.
- We also use [golangci-lint](https://golangci-lint.run/usage/install/) for linting `go` code.
- We have a pre-commit hook to enforce style/formatting.

## Next steps / Nice-to-haves:

- Make a FE app that can interact with the BE
- Don't use sqlite
- Dockerize the server
  - `Tilt` setup for local dev environments with hot-reloading would be nice
- Design the APIs better with filters for reports
- Better testing suite for db interactions rather than just a single e2e test that nukes the db afterwards to get a clean state
  - Technically the e2e test can use a new db and manage the process as well if I config out the db file and use python to control the processes.

## Prompt:

Imagine you're tasked with creating a practical solution for a local shelter to manage their donation inventory. This shelter is in need of a user-friendly tool to accurately record and track the inflow and outflow of donations, and to generate insightful reports about their donation management.

Here are the core functionalities your solution should address:

1. Donation Registration: A feature that allows the shelter staff to record details of the donations, such as the donor's name, type of donation (money, food, clothing, etc.), quantity or amount donated, and the date of the donation.
1. Donation Distribution: A feature to log when and how much of the donations are distributed, capturing the type of donation, quantity or amount distributed, and the date of distribution.
1. Donation Reports: Your solution should have the capacity to generate two types of reports: (1) An inventory report displaying the current status of donations, grouped by type. (2) A donator report, summarizing the total contributions received from each donor.
   Feel free to select any programming language that you're most comfortable with. We want you to showcase your coding abilities, problem-solving prowess, and ability to operate under ambiguity.

To clarify, you don't need to spend your time on setting up a database or creating a complex frontend. Your primary focus should be on the core functionalities mentioned above. However, if there's a particular area in which you excel, such as data modeling, frontend design, testing, or anything else, feel free to include that in your solution to show off your strengths. This will allow us to better understand if you have the abilities we currently need.

Along with your code + comments, we'd appreciate clear instructions on how to run the code.

Remember, this exercise is not about checking boxes, but about allowing you the freedom to demonstrate your unique strengths and approach to problem-solving. We're excited to see what you can create!
