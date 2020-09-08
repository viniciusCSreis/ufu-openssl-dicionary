import itertools


def main():
    base_dictionary = open("base-dictionary.txt", "r")
    words = base_dictionary.read().split("\n")

    dictionary = set({})

    modifiers = generate_modifiers()

    for word in words:
        for i in itertools.product([False, True], repeat=len(modifiers)):
            new_word = word
            for j in range(len(modifiers)):
                if i[j]:
                    new_word = modifiers[j](new_word)

            dictionary.add(new_word + "\n")

    final_dictionary = open("dictionary.txt", "w")
    for w in dictionary:
        if w != "":
            final_dictionary.write(w)

    final_dictionary.close()
    base_dictionary.close()


def generate_modifiers():
    upper = (lambda w: w.upper())
    lower = (lambda w: w.lower())
    capitalize = (lambda w: w.capitalize())
    replace_dash = (lambda w: w.replace("-", " "))
    replace_spaces = (lambda w: w.replace(" ", ""))
    replace_a = (lambda w: w.replace("a", "@").replace("A", "@"))
    replace_e = (lambda w: w.replace("e", "3").replace("E", "3"))
    replace_i = (lambda w: w.replace("i", "!").replace("I", "!"))
    to_pascal_case = (lambda word: ''.join(w for w in word.title() if not w.isspace()))
    modifiers = [upper, lower, capitalize, replace_dash, replace_spaces, replace_a, replace_e, replace_i, to_pascal_case]
    return modifiers


if __name__ == "__main__":
    main()
