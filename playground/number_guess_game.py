import os
import sys

if __name__ == "__main__":
    print("###Hello Martin! Welcome to number guessing game!")
    print("###In this game, you will give an INTEGER as your guess.")
    print("###I will give you a hint when your guess is too high or too low.")
    print("###You have 7 tries :) And you can hit ENTER on keyboard to terminate the game.")
    answer = 78
    tries = 7
    guess = 0
    finished = 0
    result = ''
    for num in range(tries):
        input_str = input()
        input_list = input_str.split('.')
        if len(input_list) > 2:
            print("invalid input: {}!".format(input_str))
            #exit(1)
        elif len(input_list) == 2:
            if input_list[0].isnumeric() and input_list[1].isnumeric():
                guess = float(input_str)
            else:
                print("invalid input: {}!".format(input_str))
                #exit(2)
        elif len(input_list) == 1:
            if input_list[0].isnumeric():
                guess = int(input_str)
            elif input_list[0] == '':
                print("The game is terminated!")
                exit(0)
            else:
                print("invalid input: {}!".format(input_str))
                #exit(3)
        else:
            print("odd input!")
            #exit(4)

        if guess < answer:
            result = "Oops!  Your guess was too low."
        elif guess > answer:
            result = "Oops!  Your guess was too high."
        elif guess == answer:
            result = "Nice!  Your guess matched the answer!"
            finished = 1
        print(result)

        if finished:
            exit(0)

    print("\n\n Oops! you have failed the game, HA HA HA, try next time~")